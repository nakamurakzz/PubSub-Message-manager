.PHONY: build run
build:
	go build -o pubusub-subscriber .

.PHONY: run
run:
	air ${project_id}

.PHONY: create-topic
create-topic:
	gcloud pubsub topics create ${topic_id}

.PHONY: create-subscription
create-subscription:
	gcloud pubsub subscriptions create ${subscription_id} --topic=${topic_id}

.PHONY: delete-topic
delete-topic:
	gcloud pubsub topics delete ${topic_id}

.PHONY: delete-subscription
delete-subscription:
	gcloud pubsub subscriptions delete ${subscription_id}