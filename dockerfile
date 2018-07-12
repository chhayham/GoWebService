# amzn linux os
FROM amazonlinux:latest

# set require env variables
ENV ENV_ROOT /usr/bin/goweb

# setup container using stacked commands
RUN echo ${ENV_ROOT} \
&& echo $PATH \
&& mkdir -p /usr/bin/goweb \
&& yum install -y tree

#VOLUME persistant storage

# change working dir
WORKDIR $ENV_ROOT

# go app listens on 8080, so only expose that port
EXPOSE 80 8080

# start app
CMD ["tree", "/usr/bin"]
