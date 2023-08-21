FROM adrienaury/go-devcontainer:v3.1

USER root

RUN apk add --update --progress --no-cache make gomplate

ARG VERSION_MILLER=6.2.0
RUN    wget -nv -O- https://github.com/johnkerl/miller/releases/download/v${VERSION_MILLER}/miller-${VERSION_MILLER}-linux-amd64.tar.gz | tar xz --strip-components 1 -C /usr/bin miller-${VERSION_MILLER}-linux-amd64/mlr \
    && chmod +x /usr/bin/mlr

ARG VERSION_PIMO=1.19.0
RUN wget -O- https://github.com/CGI-FR/PIMO/releases/download/v${VERSION_PIMO}/pimo_${VERSION_PIMO}_linux_amd64.tar.gz | tar xz -C /usr/bin pimo

USER vscode
