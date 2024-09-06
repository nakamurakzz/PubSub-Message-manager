# Google Cloud Pub/Sub Reciever

This repository contains a simple web application that allows users to subscribe to Google Cloud Pub/Sub topics and receive messages in real-time.

## Features

- **Project ID Registration:** Users can register their Google Cloud Project ID to access Pub/Sub topics.
- **Topic Selection:** Users can select a topic from a list of available topics.
- **Real-Time Message Display:** Messages published to the selected topic are displayed in real-time on the web page.

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/pubsub-interface.git
   cd pubsub-interface
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up your Google Cloud project:**
   - Create a new Google Cloud project or use an existing one.
   - Enable the Pub/Sub API in your project.
   - Create a Pub/Sub topic.
   ```bash
   make create-topic project_id={your_project_id} topic_id={your_topic_id}
   ```
   - Create a Pub/Sub subscription for the topic.
   ```bash
   make create-subscription project_id={your_project_id} topic_id={your_topic_id} subscription_id={your_subscription_id}
   ```

4. **Run the application:**
   ```bash
   make build
   make run project_id={your_project_id} port={your_port}
   ```
   Replace `<your_project_id>` with your actual Google Cloud Project ID.
   `<your_port>` is optional. If you don't specify it, the default port is 5000.

5. **Access the application:**
   Open your web browser and navigate to `http://localhost:{your_port}`.

## Usage

1. Select a topic from the dropdown list.
2. Click the "Subscribe" button.
3. Messages published to the selected topic will be displayed in real-time.

## Makefile

The `Makefile` provides commands for building, running, and managing Pub/Sub resources:

- **`build`:** Builds the application.
- **`run`:** Runs the application.
- **`run-dev`:** Runs the application with hot reload.
- **`create-topic`:** Creates a new Pub/Sub topic.
- **`create-subscription`:** Creates a new Pub/Sub subscription for a topic.
- **`delete-topic`:** Deletes a Pub/Sub topic.
- **`delete-subscription`:** Deletes a Pub/Sub subscription.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.