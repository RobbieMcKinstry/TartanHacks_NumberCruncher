from flask import Flask
from flask import request as flask_req
import requests as r

app = Flask(__name__)

@app.route("/", methods=['GET'])
def hello():
    """This returns the index.html page"""
    return "Hello World!"

"""Make sure you're getting a GET request with the url including ?symbol="""
@app.route('/new_stock', methods=['GET'])
def new_stock_request():
    """This endpoint receives the json, makes the request to the Go server, and returns the results"""
    symbol = flask_req.args.get('symbol') # gives you the symbol, as a string

if __name__ == "__main__":
    app.debug = True
    app.run(host='0.0.0.0')
