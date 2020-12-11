FROM library/golang:1.11

RUN mkdir -p /root/go

COPY "docker/init.sh" "/data/init.sh"

RUN mkdir -p /root/go/src/github.com/yaoice/istio-crd-tool/

COPY . /root/go/src/github.com/yaoice/istio-crd-tool/

RUN /bin/bash -c 'echo "export GOPATH=\"/root/go\"" >> /root/.bashrc \
    && source /root/.bashrc \
    && cd $GOPATH \
    && mkdir -p /etc/istio-crd-tool \
    && cp $GOPATH/src/github.com/yaoice/istio-crd-tool/etc/istio-crd-tool.toml /etc/istio-crd-tool/istio-crd-tool.toml \
    && go build -o $GOPATH/bin/ice github.com/yaoice/istio-crd-tool/cmd \
    && rm -rf $GOPATH/src/github.com/yaoice/istio-crd-tool \
    && chmod +x /data/init.sh'

EXPOSE 9094

CMD ["/data/init.sh"]