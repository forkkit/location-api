FROM scratch
ADD geo-api /
ENTRYPOINT [ "/geo-api" ]
