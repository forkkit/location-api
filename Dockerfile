FROM scratch
ADD location-api /
ENTRYPOINT [ "/location-api" ]
