from src import UniversalModel
from sanic import Sanic, response
from ray import serve

# Initialize Sanic app
app = Sanic("Catwalk")

# Initialize Ray and Ray Serve
serve.start(detached=True)

UniversalModel.options(name="tf_model").deploy("model_keras.h5")
UniversalModel.options(name="torch_model").deploy("model_torch.pt")

# Sanic endpoint to forward requests to Ray Serve
@app.post("/<model_name>")
async def serve_model(request, model_name):
    # Forward the request to the correct Ray Serve model
    serve_handle = serve.get_deployment(model_name).get_handle()
    prediction = await serve_handle.remote(data=request.json)
    return response.json(prediction)

if __name__ == "__main__":
    # Run the Sanic app
    app.run(host="0.0.0.0", port=8000)
