const { UserServiceClient } = require("./user_grpc_web_pb");
const { GetUserRequest } = require("./user_pb");

const client = new UserServiceClient("http://localhost:8081"); // Envoy proxy

document.getElementById("btn").addEventListener("click", () => {
  const req = new GetUserRequest();
  req.setId(123); // harus number, bukan string

  client.getUser(req, {}, (err, response) => {
    if (err) {
      document.getElementById("output").innerText = "Error: " + err.message;
    } else {
      document.getElementById("output").innerText = JSON.stringify(
        response.toObject(),
        null,
        2
      );
    }
  });
});
