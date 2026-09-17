#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

INSTALL_BIN="/usr/local/bin/ldap-extractor"

CONFIG_DIR="/etc/ldap-extractor"
CONFIG_FILE="${CONFIG_DIR}/config.yaml"

DATA_DIR="/opt/ldap-extractor"
FILTER_DIR="${DATA_DIR}/filtered"
DIFF_DIR="${DATA_DIR}/diff"
SYSTEM_SERVICE="/etc/systemd/system"
SERVICE_NAME="ldap-extractor"

usage() {
    cat <<EOF
Usage:
  $0 --install
  $0 --uninstall

Options:
  --install      Install the latest ldap-extractor binary
  --uninstall    Remove ldap-extractor
  --help         Show this help

Examples:
  sudo $0 --install
  sudo $0 --uninstall
EOF
}


require_root() {
    if [[ "${EUID}" -ne 0 ]]; then
        echo "Error: this operation must be run as root."
        echo "Use sudo."
        exit 1
    fi
}


find_binary() {

    local binary

    binary="$(
        find "${SCRIPT_DIR}" \
            -maxdepth 1 \
            -type f \
            -name 'ldap-extractor-v*' \
            -perm -u+x \
            -print \
        | sort -V \
        | tail -n 1
    )"

    if [[ -z "${binary}" ]]; then
        echo "Error: no ldap-extractor binary found."
        echo
        echo "Expected a binary such as:"
        echo "  ldap-extractor-v0.1.2"
        exit 1
    fi

    echo "${binary}"
}


install_ldap_extractor() {

    local binary
    
    binary="$(find_binary)"

    echo "Found binary:"
    echo "  ${binary}"

    echo
    echo "Installing as:"
    echo "  ${INSTALL_BIN}"


    # Install the versioned binary using the stable command name.
    install -Dm755 \
        "${binary}" \
        "${INSTALL_BIN}"
    
    install -Dm644 \
        "${SERVICE_NAME}.service" \
        "${SYSTEM_SERVICE}"
    
    install -Dm644 \
        "${SERVICE_NAME}.timer" \
        "${SYSTEM_SERVICE}"


    # --------------------------------------------------------------
    # Configuration
    # --------------------------------------------------------------

    mkdir -p "${CONFIG_DIR}"


    if [[ ! -f "${CONFIG_FILE}" ]]; then

        cat > "${CONFIG_FILE}" <<'EOF'
ldap:
  url: "ldap://<ip>:389"
  username: "username"
  password: "passwd"
  base_dn: "dc=example,dc=com"

  search:
    filter: "(objectCategory=person)"
    page_size: 3000
    attributes:
      - sAMAccountName
      - department
      - employeeID
      - cn

  ldif:
    encode_non_ascii: true

output:
  file: "/opt/ldap-extractor/ldap_dump.ldif"

  filter:
    directory: "/opt/ldap-extractor/filtered"
    suffix: "filtered_users.json"
    retention: 5

  diff:
    directory: "/opt/ldap-extractor/diff"
    suffix: "diff_users.json"
    retention: 5
EOF

        chmod 600 "${CONFIG_FILE}"

        echo
        echo "Created configuration:"
        echo "  ${CONFIG_FILE}"

    else

        echo
        echo "Configuration already exists:"
        echo "  ${CONFIG_FILE}"

        echo "Keeping existing configuration."
    fi


    # --------------------------------------------------------------
    # Output directories
    # --------------------------------------------------------------

    mkdir -p \
        "${DATA_DIR}" \
        "${FILTER_DIR}" \
        "${DIFF_DIR}"


    chmod 755 \
        "${DATA_DIR}" \
        "${FILTER_DIR}" \
        "${DIFF_DIR}"
    
    systemctl enable --now ${SERVICE_NAME}.timer
    systemctl daemon-reload


    echo
    echo "Output directories:"
    echo "  ${DATA_DIR}"
    echo "  ${FILTER_DIR}"
    echo "  ${DIFF_DIR}"


    echo
    echo "Installation complete."
    echo
    echo "IMPORTANT:"
    echo "  Configure your LDAP credentials and connection settings in:"
    echo "    ${CONFIG_FILE}"
    echo
    echo "  In particular, update:"
    echo "    - ldap.url"
    echo "    - ldap.username"
    echo "    - ldap.password"
    echo "    - ldap.base_dn"
    echo 
    echo "  Configure the ldap-extractor.timer based on how often the extractor has to run."
    echo
    echo "Run:"
    echo "  ldap-extractor --help"
}


uninstall_ldap_extractor() {

    service_path="/etc/systemd/system"

    echo "Uninstalling ldap-extractor..."


    if [[ -f "${INSTALL_BIN}" ]]; then

        rm -f "${INSTALL_BIN}"

        echo "Removed:"
        echo "  ${INSTALL_BIN}"

    else

        echo "Binary is not installed:"
        echo "  ${INSTALL_BIN}"

    fi

    if [[ -f "${service_path}/${SERVICE_NAME}.service" ]]; then

        rm -f "${service_path}/${SERVICE_NAME}.service"

        echo "Removed:"
        echo "  ${service_path}/${SERVICE_NAME}.service"

    else

        echo "${SERVICE_NAME}.service is not installed:"
        echo "  ${service_path}/${SERVICE_NAME}.service"
    fi

    if [[ -f "${service_path}/${SERVICE_NAME}.timer" ]]; then

        rm -f "${service_path}/${SERVICE_NAME}.timer"

        echo "Removed:"
        echo "  ${service_path}/${SERVICE_NAME}.timer"

    else

        echo "${SERVICE_NAME}.timer is not installed:"
        echo "  ${service_path}/${SERVICE_NAME}.timer"
    fi

    echo
    echo "Uninstallation complete."
    echo
    echo "Configuration and application data were preserved:"
    echo "  ${CONFIG_FILE}"
    echo "  ${DATA_DIR}"
}


# --------------------------------------------------------------
# Main
# --------------------------------------------------------------

if [[ $# -ne 1 ]]; then
    usage
    exit 1
fi


case "$1" in

    --install)
        require_root
        install_ldap_extractor
        ;;


    --uninstall)
        require_root
        uninstall_ldap_extractor
        ;;


    --help|-h)
        usage
        ;;


    *)
        echo "Error: unknown option: $1"
        echo
        usage
        exit 1
        ;;

esac
