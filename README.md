# TripleaWire : Capture your network traffic and analyze it.

TripleaWire is a network traffic capture and analysis tool. It is designed to capture network traffic and analyze it in real-time. The tool is built using golang for the core and React for the UI.

## Features
- Capture network traffic in real-time (Only TCP and ICMP packets are supported for now.)
- Show packets details in a human readable format.
- Trigger alerts based on packet data.
- Websocket server to stream packets to the UI.

## Here is an sample example of the UI:

![TripleaWire UI](https://github.com/kraaakilo/tripleawire/)

### Running the Core Project:
1. **Install Go**: Download and install Go from [here](https://golang.org/dl/).

2. **Clone Repository**: Clone the TripleaWire repository.
    ```bash
    git clone https://github.com/kraaakilo/tripleawire.git
    ```

3. **Build Core Project**: Navigate to the core directory and run the following command to build the core project.
    ```bash
    cd tripleawire/core
    go build
    ```

4. **Start Websocket Server**: Run the following command to start the core websocket server.
    ```bash
    ./triplewire --interface interface-to-use --mode web
    ```
    Optionally, you can run the following command to start the core in CLI mode.
    ```bash
    ./triplewire --interface interface-to-use --mode cli
    ```

### Running the UI:
1. **Install Node and pnpm**: Ensure you have Node.js and pnpm installed on your machine.

2. **Clone Repository**: Clone the TripleaWire repository.
    ```bash
    git clone https://github.com/kraaakilo/tripleawire.git
    ```

3. **Install Dependencies**: Navigate to the UI directory and install dependencies.
    ```bash
    cd tripleawire/ui
    pnpm install
    ```

4. **Start UI**: Run the following command to start the UI.
    ```bash
    pnpm dev
    ```

### Important Note:
- **Caution**: Since the project is still in development, it's advisable to use it with caution.
- **Administrator Privileges**: Run the program as an administrator to allow network traffic capture.

The provided guide should help users get started with setting up and running TripleaWire for network traffic analysis.