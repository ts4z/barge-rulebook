#
# This Dockerfile can be used to make a docker image that can build the
# LaTeX-based pdf version of the rulebook.  Unfortunately I have been unable to
# use that Docker image to automate producing the rulebook PDF, as it is over a
# gigabyte in size due to the excessive number of fonts.
#

FROM debian:12-slim
RUN apt update && apt upgrade && apt install -y texlive-xetex texlive-fonts-extra libcommonmark-perl make
RUN mkdir /work
RUN chown nobody /work
WORKDIR /work
#USER nobody

CMD ["bash"]
