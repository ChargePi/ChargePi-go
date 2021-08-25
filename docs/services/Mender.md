## ChargePi with Mender

We will use [Mender](https://mender.io/) as OTA service for updating Linux, Docker and the client. For the OCPP
server/central system, we will use **[SteVe](https://github.com/RWTH-i5-IDSG/steve)**.

1. Sign up on Mender or install your
   own [production Mender server](https://docs.mender.io/3.0/server-installation/installation-with-docker-compose).
2. Follow this guide
   for [installing Mender on the Pi](https://docs.mender.io/get-started/preparation/prepare-a-raspberry-pi-device).
3. Add the Raspberry to the Devices list and test if Mender works.
5. Clone the SteVe repository to ChargePi/. It should automatically create steve folder.
   *Skip this step if you do not want the server to be hosted on the same machine. Mostly used for testing purposes.*
6. If you do not wish to automatically update the OCPP client, change this label to false in **docker-compose.yaml**:

```
     chargepi:
        label:
           - com.centurylinklabs.watchtower.enable="false"
```

   **Not recommended when using Mender, since you can update the images with Mender.**

7. Run docker-compose:

   > docker-compose up -d

* Optionally, you can run only the client service:

  > docker-compose up -d chargepi

#### Helpful references

- [Installing Docker on Pi](https://www.docker.com/blog/happy-pi-day-docker-raspberry-pi/)
- [Watchtower](https://github.com/containrrr/watchtower)
- [Docker](https://docs.docker.com/)
- [Pinout](https://pinout.xyz/)
- [Mender docs](https://docs.mender.io/get-started/preparation/prepare-a-raspberry-pi-device)