#
# Install dependancies and run install script
#

FROM ubuntu:22.04

ADD ["install.sh", "frames.sh", "convert.sh", "run.sh", "scripts/"]
RUN chmod -R +x scripts/
RUN scripts/install.sh

#
# Execute via make_thumbnail command
#

CMD ["scripts/make_thumbnail.sh", "sample.mp4", "output.gif"]
