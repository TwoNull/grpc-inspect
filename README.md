<h3 align="center">gRPC Inspect</h3>
  <p align="center">
    A Go implementation of a gRPC client and server to pull item inspect data from the Counter-Strike: Global Offensive (CSGO) game coordinator.
    <br />
    <br />
    <a href="https://nullptrs.co/" target="_blank" rel="noopener noreferrer"><strong>Read the Full Writeup »</strong></a>
    <br />
  </p>
</div>

### Disclaimer
 This server is not optimized for a production environment, as certain race conditions can occur under high load and valid requests may be dropped. Use at your own risk.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Requirements
 Below are the prerequisites to run this tool.
 * [Go ^1.20](https://go.dev/dl/)

### Setup
 1. Clone the repo
   ```
   $ git clone https://github.com/twonull/grpc-inspect.git
   ```
 2. Navigate to the project directory and install dependencies
   ```
   $ cd grpc-inspect
   $ go get .
   ```
 3. Create an accounts.txt file with Steam accounts in the following format. 
    Note that these accounts **must** have Steam Guard/Email 2FA disabled
   ```
   username1:password1
   username2:password2
   ```
 4. Start the gRPC server. Replace `ACCOUNTS` with the path of the accounts.txt file you just created
   ```
   $ GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore go run ./server -file ACCOUNTS
   ```
 5. Run the example client script to receive a response from the server
   ```
   $ go run ./client
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Made With

[![Golang][Golang]][Go-url]
[![gRPC][gRPC]][gRPC-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- LICENSE -->
### License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
[Golang]: https://shields.io/badge/Golang-5DC9E2?style=for-the-badge&logo=Go&logoColor=FFF
[Go-url]: https://go.dev/
[gRPC]: https://shields.io/badge/gRPC-44969D?style=for-the-badge&logo=&logoColor=FFF
[gRPC-url]: https://grpc.io/