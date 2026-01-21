# Docker file builds an image with a tinode chat server.
#
# In order to run the image you have to link it to a running database container. For example, to
# to use RethinkDB (named 'rethinkdb') and map the port where the tinode server accepts connections:
#
# $ docker run -p 6060:6060 -d --link rethinkdb \
#	--env UID_ENCRYPTION_KEY=base64+encoded+16+bytes= \
#	--env API_KEY_SALT=base64+encoded+32+bytes \
#	--env AUTH_TOKEN_KEY=base64+encoded+32+bytes \
#	tinode-server

FROM alpine:3.22

ARG VERSION=0.25
ENV VERSION=$VERSION
ARG BINVERS=$VERSION

LABEL maintainer="Tinode Team <info@tinode.co>"
LABEL name="TinodeChatServer"
LABEL version=$VERSION

# Build-time options.

# Database selector. Builds for MySQL by default.
# Alternatively use one of: postgres mongodb rethinkdb for a corresponsing
# DB backend or alldbs to build a generic Tinode docker image, for example:
# `--build-arg TARGET_DB=postgres` to build for PostgreSQL.
ARG TARGET_DB=mysql
ENV TARGET_DB=$TARGET_DB

# Runtime options.

# Specifies the database host:port pair to wait for before running Tinode.
# Ignored if empty.
ENV WAIT_FOR=

# An option to reset database.
ENV RESET_DB=false

# An option to upgrade database.
ENV UPGRADE_DB=false

# Option to skip DB initialization when it's missing.
ENV NO_DB_INIT=false

# Load sample data to database from data.json.
ARG SAMPLE_DATA=data.json
ENV SAMPLE_DATA=$SAMPLE_DATA

# Default country code to use in communication.
ENV DEFAULT_COUNTRY_CODE=US

# The MySQL DSN connection.
ENV MYSQL_DSN='root:123456@tcp(mysql-server:3306)/tinode?parseTime=true&collation=utf8mb4_0900_ai_ci'

# The PostgreSQL DSN connection.
ENV POSTGRES_DSN='postgresql://postgres:123456@postgres-server:5432/tinode?sslmode=disable&connect_timeout=10'

# Disable chatbot plugin by default.
ENV PLUGIN_PYTHON_CHAT_BOT_ENABLED=false

# Default handler for large files
ENV MEDIA_HANDLER=fs

# Whitelisted domains for file and S3 large media handler.
ENV FS_CORS_ORIGINS='["*"]'
ENV AWS_CORS_ORIGINS='["*"]'

# AWS S3 parameters
ENV AWS_ACCESS_KEY_ID=
ENV AWS_SECRET_ACCESS_KEY=
ENV AWS_REGION=
ENV AWS_S3_BUCKET=
ENV AWS_S3_ENDPOINT=

# Default externally-visible hostname for email verification.
ENV SMTP_HOST_URL='http://localhost:6060'
# Email parameters decalarations.
ENV SMTP_SERVER=
ENV SMTP_PORT=
ENV SMTP_SENDER=
ENV SMTP_LOGIN=
ENV SMTP_PASSWORD=
ENV SMTP_AUTH_MECHANISM=
ENV SMTP_HELO_HOST=
ENV EMAIL_VERIFICATION_REQUIRED=
ENV DEBUG_EMAIL_VERIFICATION_CODE=

# Whitelist of permitted email domains for email verification (empty list means all domains are permitted)
ENV SMTP_DOMAINS=''

# Various encryption and salt keys. Replace with your own in production.

# Salt used to generate the API key. Don't change it unless you also change the
# API key in the webapp & Android.
ENV API_KEY_SALT=T713/rYYgW7g4m3vG6zGRh7+FM1t0T8j13koXScOAj4=

# Key used to sign authentication tokens.
ENV AUTH_TOKEN_KEY=wfaY2RgF2S1OQI/ZlK+LSrp1KB2jwAdGAIHQ7JZn+Kc=

# Key to initialize UID generator
ENV UID_ENCRYPTION_KEY=la6YsO+bNX/+XIkOqc5Svw==

# Disable TLS by default.
ENV TLS_ENABLED=false
ENV TLS_DOMAIN_NAME=
ENV TLS_CONTACT_ADDRESS=

# Disable push notifications by default.
ENV FCM_PUSH_ENABLED=false
# Declare FCM-related vars
ENV FCM_API_KEY=
ENV FCM_APP_ID=
ENV FCM_SENDER_ID=
ENV FCM_PROJECT_ID=
ENV FCM_VAPID_KEY=
ENV FCM_MEASUREMENT_ID=

# Enable Android-specific notifications by default.
ENV FCM_INCLUDE_ANDROID_NOTIFICATION=true

# Disable push notifications via Tinode Push Gateway.
ENV TNPG_PUSH_ENABLED=false

# Tinode Push Gateway authentication token.
ENV TNPG_AUTH_TOKEN=

# Tinode Push Gateway organization name as registered at console.tinode.co
ENV TNPG_ORG=

# Video calls configuration.
ENV WEBRTC_ENABLED=false
ENV ICE_SERVERS_FILE=

# Use the target db by default.
# When TARGET_DB is "alldbs", it is the user's responsibility
# to set STORE_USE_ADAPTER to the desired db adapter correctly.
ENV STORE_USE_ADAPTER=$TARGET_DB

# Url path for exposing the server's internal status. E.g. '/status'
ENV SERVER_STATUS_PATH=''

# Garbage collection of unfinished account registrations.
ENV ACC_GC_ENABLED=false

# Install root certificates, they are needed for email validator to work
# with the TLS SMTP servers like Gmail or Mailjet. Also add bash and grep.
# Add build dependencies: go, git, make, gcc, musl-dev
RUN apk update && \
	apk add --no-cache ca-certificates bash grep go git make gcc musl-dev

WORKDIR /opt/tinode

# Copy the entire project source code
COPY . .

# Build the server and init-db
RUN cd server && \
	if [ "$TARGET_DB" = "alldbs" ]; then \
		buildtag="mysql mongodb rethinkdb postgres"; \
	else \
		buildtag="$TARGET_DB"; \
	fi && \
	go build -ldflags "-s -w" -tags "$buildtag" -o tinode . && \
	cp tinode /opt/tinode/ && \
	cp tinode.conf /opt/tinode/ && \
	cp -r templ /opt/tinode/ && \
	cd ../tinode-db && \
	go build -ldflags "-s -w" -tags "$buildtag" -o init-db . && \
	cp init-db /opt/tinode/ && \
	cp credentials.sh /opt/tinode/ && \
	cp data.json /opt/tinode/ && \
	cp *.jpg /opt/tinode/

# Copy config template and entrypoint script
COPY docker/tinode/config.template .
COPY docker/tinode/entrypoint.sh .

# Clean up build dependencies and source code
RUN apk del go git make gcc musl-dev && \
	rm -rf server tinode-db keygen loadtest monitoring pbx py_grpc rest-auth tn-cli

# Create directory for chatbot data and static files.
RUN mkdir /botdata && \
	mkdir -p static

# Make scripts runnable and fix line endings
RUN apk add --no-cache dos2unix && \
	dos2unix entrypoint.sh && \
	dos2unix credentials.sh && \
	chmod +x entrypoint.sh && \
	chmod +x credentials.sh && \
	apk del dos2unix

# Healthcheck: check if port 6060 is open
HEALTHCHECK --interval=1m --timeout=3s --start-period=30s \
  CMD nc -z localhost 6060 || exit 1

# Generate config from template and run the server.
ENTRYPOINT ["./entrypoint.sh"]

# HTTP, gRPC, cluster ports
EXPOSE 6060 16060 12000-12003
