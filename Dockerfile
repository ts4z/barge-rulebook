
FROM debian:12-slim
RUN apt update && apt upgrade && apt install -y texlive-xetex libcommonmark-perl make texlive-fonts-extra
USER nobody
