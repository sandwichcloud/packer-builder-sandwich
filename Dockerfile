FROM hashicorp/packer:1.1.3

COPY bin/packer-builder-sandwich_linux_amd64 /root/.packer.d/plugins/packer-builder-sandwich