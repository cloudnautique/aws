'use strict';

const axios = require("axios").default;

exports.handler = (event: { queryStringParameters: { url: string; }; }) => {
    if (!event.queryStringParameters || !event.queryStringParameters.url) {
        let response = {
            statusCode: 200,
            headers: {
                "content-type": "text/plain; charset=utf-8"
            },
            body: "Please provide a url as a query string parameter"
        };
        return (Promise.resolve(response));
    } else {
        const url = event.queryStringParameters.url;
        return call(url)
    }
};


function call(url: string) {
    console.log("Getting: " + url);
    return axios
        .get(url)
        .then((response: { data: string }) => {
            console.log("Got content for request: " + url);
            return {
                statusCode: 200,
                headers: {
                    "content-type": "text/plain; charset=utf-8"
                },
                body: response.data
            };
        })
        .catch((error: Error) => {
            return {
                statusCode: 500,
                headers: {
                    "content-type": "text/plain; charset=utf-8"
                },
                body: "Some error fetching the content\n" + error + "\n"
            };
        });
}