FROM moby/buildkit:v0.9.3
WORKDIR /web
COPY web README.md /web/
ENV PATH=/web:$PATH
ENTRYPOINT [ "/bhojpur/web" ]