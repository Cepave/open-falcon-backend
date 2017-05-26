#!/bin/bash
cd /go/src/"$GET_NAME" && source docker/build-alpine/scripts/logs.sh 

## -------------------------------------------- ##
## Upload binaries

UPLOAD_S3() {
	if [ "$AWS_KEY" != "" ] && [ "$AWS_SECRET" != "" ] ; then
		log "s3 syncing..."
		aws s3 sync bin "$AWS_S3_BUCKET"/"$GIT_BRANCH" --storage-class REDUCED_REDUNDANCY --acl public-read
	fi
}

UPLOAD_FTP() {
	NEW_BINS=$(find /go/src/$GET_NAME/bin -mmin -3 -type f)
	if [ "$NEW_BINS" != "" ] && [ "$FTP_AUTH" != "" ] && [ "$FTP_URL" != "" ] ; then
		for BIN in $NEW_BINS
		do
			log "$BIN is uploading to FTP" 
			curl -s -F "file=@$BIN" -u "$FTP_AUTH" "$FTP_URL" 2>&1
		done
		log "Upload finished"
	fi
}

## -------------------------------------------- ##
## Go Get & Build

MAKEALL(){
	log "Building..." && make all
	cp -f open-falcon bin
	find /go/src/$GET_NAME/bin -type f | xargs ls -al | awk '{print $5, $6, $7, $8, $9}' > /go/src/$GET_NAME/bin/info.txt
}

BUILD() {
	cd /go/src && go get "$GET_NAME" && cd /go/src/$GET_NAME

	git checkout $GIT_BRANCH && if [ "$GIT_TAG" != "" ]; then git checkout "$GIT_TAG"; fi
	MAKEALL && UPLOAD_S3 && UPLOAD_FTP

	while [ "$GIT_TAG" == "" ]
	do 
		if [ "$(git pull | grep up-to-date)" == "" ]; then
			log "Go get codes..." && go get "$GET_NAME"
			MAKEALL && UPLOAD_S3 && UPLOAD_FTP
		fi

		sleep "$INTERVAL"
	done
}

## -------------------------------------------- ##
## Actions

BUILD
