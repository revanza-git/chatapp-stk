# 🚀 Google Colab + Ollama Setup Guide

## Step 1: Setup Ollama in Google Colab

1. **Open Google Colab**: Go to [colab.research.google.com](https://colab.research.google.com)

2. **Upload the notebook**: Upload the `colab_ollama_setup.ipynb` file

3. **Run all cells** in sequence:

   - Cell 1: Installs Ollama
   - Cell 2: Starts Ollama service
   - Cell 3: Downloads Llama 3.1 8B model (~5-10 minutes)
   - Cell 4: Tests the model locally
   - Cell 5: Sets up ngrok tunnel (gives you public URL)
   - Cell 6: Tests the public API
   - Cell 7: Keeps service running

4. **Copy the ngrok URL** from Step 5 output (looks like `https://xyz123.ngrok.io`)

## Step 2: Configure Your Go Backend

1. **Set the Ollama URL**:

   ```powershell
   $env:OLLAMA_URL="https://xyz123.ngrok.io"  # Replace with your actual URL
   ```

2. **Keep your HF token** (as backup):

   ```powershell
   $env:HF_TOKEN="hf_your_actual_token_here"
   ```

3. **Start your Go server**:
   ```powershell
   go run main.go
   ```

## Step 3: Test Your Setup

Your Go backend will now use APIs in this priority order:

1. **Ollama API** (from Google Colab) - Free, powerful
2. **Hugging Face API** (your token) - Free, limited
3. **Mock responses** - Always available

### Test API Priority:

```bash
# Test with curl
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is our password policy?",
    "type": "onboarding"
  }'
```

## Step 4: Expected Output

### Console Output:

```
🚀 Security Chatbot Server starting on :8080...
📝 Configuration:
   ✅ Ollama URL: https://xyz123.ngrok.io
   ✅ Hugging Face token configured
   ℹ️  Mock responses available as fallback

✅ Using Ollama API from Google Colab
```

### API Response:

```json
{
  "response": "Our password policy requires passwords to be at least 12 characters long with a combination of uppercase letters, lowercase letters, numbers, and special characters. Passwords must be changed every 90 days to maintain security. Additionally, passwords should not contain personal information and should be unique across different systems.",
  "type": "onboarding"
}
```

## Benefits of This Setup:

✅ **Completely Free** - Google Colab provides free GPU/RAM  
✅ **Powerful Model** - Llama 3.1 8B is much better than smaller models  
✅ **No Local Resources** - Runs entirely in the cloud  
✅ **Easy Testing** - Perfect for development and prototyping  
✅ **Automatic Fallback** - Falls back to HF API or mocks if Colab stops

## Tips:

- **Keep Colab running**: The ngrok tunnel only works while Colab is active
- **Session limits**: Google Colab has usage limits (usually 12 hours)
- **Restart process**: If Colab disconnects, just re-run the cells
- **Production ready**: For production, consider deploying to GCP/AWS

## Troubleshooting:

### If Ollama URL doesn't work:

1. Check if Colab is still running
2. Verify the ngrok URL is correct
3. Check firewall settings

### If no responses:

1. Check console logs for which API is being used
2. Verify environment variables are set
3. Test individual API endpoints

### Model download fails:

1. Restart Colab session
2. Try a smaller model: `ollama pull llama3.1:3b`
3. Check Colab's disk space
