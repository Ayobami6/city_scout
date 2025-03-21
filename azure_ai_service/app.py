from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import openai
import base64
import os
import traceback
from openai import AzureOpenAI
# from openai import AzureOpenAI
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize FastAPI app
app = FastAPI()

class Query(BaseModel):
    prompt: str



@app.post("/safe_route")
async def get_safe_route(query: Query):
    try:
        endpoint = "https://ayoba-m8gwjbda-northcentralus.openai.azure.com/"  
        deployment = "gpt-4o"  
        subscription_key =  os.getenv("AZURE_OPENAI_API_KEY")

        # Initialize Azure OpenAI Service client with key-based authentication    
        client = AzureOpenAI(  
            azure_endpoint=endpoint,  
            api_key=subscription_key,  
            api_version="2024-05-01-preview",
        )
        with open("textfile.txt", "rb") as file:
            encoded_content = base64.b64encode(file.read()).decode("utf-8")
            
            
        # IMAGE_PATH = "YOUR_IMAGE_PATH"
        # encoded_image = base64.b64encode(open(IMAGE_PATH, 'rb').read()).decode('ascii')

        #Prepare the chat prompt 
        chat_prompt = [
            {
                "role": "system",
                "content": [
                    {
                        "type": "text",
                        "text": "You are an AI assistant that helps people find information., Here is a file: " + encoded_content + " for you to use as a data source"
                    }
                ]
            },
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": query.prompt
                    }
                ]
            },
        ] 
            
        # Include speech result if speech is enabled  
        messages = chat_prompt  
            
        # Generate the completion  
        completion = client.chat.completions.create(  
            model=deployment,
            messages=messages,
            max_tokens=800,  
            temperature=0.7,  
            top_p=0.95,  
            frequency_penalty=0,  
            presence_penalty=0,
            stop=None,  
            stream=False
        )

        # Extract the response text
        response_text = completion.choices[0].message.content
        return {"response": response_text}
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "healthy"}
