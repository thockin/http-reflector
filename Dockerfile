FROM scratch
COPY http-reflector /
EXPOSE 80
ENTRYPOINT ["/http-reflector"]
