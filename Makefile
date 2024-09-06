.PHONY: build
build:
	go build -o cmd/pubusub-subscriber .

.PHONY: run
run:build
	./cmd/pubusub-subscriber ${project_id} ${port}

.PHONY: run-dev
run-dev:
	air ${project_id}

.PHONY: create-topic
create-topic:
	gcloud pubsub topics create ${topic_id} --project=${project_id}

.PHONY: create-subscription
create-subscription:
	gcloud pubsub subscriptions create ${subscription_id} --topic=${topic_id} --project=${project_id}

.PHONY: delete-topic
delete-topic:
	gcloud pubsub topics delete ${topic_id} --project=${project_id}

.PHONY: delete-subscription
delete-subscription:
	gcloud pubsub subscriptions delete ${subscription_id} --project=${project_id}