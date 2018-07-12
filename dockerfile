# amzn linux os
FROM amazonlinux:latest

# set require env variables
ENV ENV_ROOT /usr/bin

# setup container using stacked commands
RUN echo ${ENV_ROOT} \
&& echo $PATH

# go app listens on 8080, so only expose that port
EXPOSE 80 8080

# start app
CMD ["uname", "-a"]
