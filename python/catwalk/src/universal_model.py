import numpy as np
import tensorflow as tf
import torch
import onnxruntime as ort
from sklearn.externals import joblib  # Use joblib to load scikit-learn models
from ray import serve

def load_model(model_path):
    if model_path.endswith('.h5'):
        return "tensorflow", tf.keras.models.load_model(model_path)
    elif model_path.endswith('.pt'):
        model = torch.load(model_path)
        model.eval()
        return "torch", model
    elif model_path.endswith('.onnx'):
        session = ort.InferenceSession(model_path)
        return "onnx", session
    elif model_path.endswith('.joblib') or model_path.endswith('.pkl'):
        model = joblib.load(model_path)
        return "sklearn", model
    else:
        raise ValueError("Unsupported model format")

@serve.deployment
class UniversalModel:
    def __init__(self, model_path):
        self.model_type, self.model = load_model(model_path)

    async def __call__(self, request):
        json_input = await request.json()
        input_data = json_input.get("data")

        # Convert input data to the appropriate format for each model type
        # TODO: This needs to be updated to handle introspect and handle different input types
        
        if self.model_type == "tensorflow":
            prediction = self.model.predict([input_data]).tolist()
        elif self.model_type == "torch":
            input_tensor = torch.tensor(input_data, dtype=torch.float32)
            with torch.no_grad():
                prediction = self.model(input_tensor).numpy().tolist()
        elif self.model_type == "onnx":
            ort_inputs = {self.model.get_inputs()[0].name: np.array(input_data, dtype=np.float32)}
            ort_outs = self.model.run(None, ort_inputs)
            prediction = ort_outs[0].tolist()
        elif self.model_type == "sklearn":
            prediction = self.model.predict([input_data]).tolist()
        
        return {"prediction": prediction}
