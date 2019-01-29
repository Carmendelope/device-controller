
# iron/go is the alpine image with only ca-certificates added
FROM iron/go

# Now just add the binary
COPY device-controller /nalej/

ENTRYPOINT ["./nalej/device-controller"]