const fs = require('fs').promises;
const https = require('https');
const http = require('http');

const filePath = 'tests.json';


function RunTest(obj) {
    return new Promise((resolve) => {
        let Ok = true;
        const client = obj.Url.startsWith('https') ? https : http;

        // Prepare the data for POST requests
        const postData = JSON.stringify(obj.Content);

        // Define the request options, including method and headers
        const options = {
            method: obj.Method,
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': Buffer.byteLength(postData) // Set the content length for POST requests
            }
        };

        const req = client.request(obj.Url, options, (res) => {
            let responseData = '';

            // Collect the response data
            res.on('data', (chunk) => {
                responseData += chunk;
            });

            // On the response end, evaluate the results
            res.on('end', () => {
                // Check if the response status is as expected
                if (res.statusCode !== obj.Expected) {
                    console.error(`Test (${obj.Name}) failed: Expected status ${obj.Expected}, but got ${res.statusCode}`);
                    Ok = false;
                }

                // Check if the response data matches expected content
                if (responseData !== obj.ExpectedContent) {
                    console.error(`Test (${obj.Name}) failed: Expected content "${obj.ExpectedContent}", but got "${responseData}"`);
                    Ok = false;
                }

                console.log("Test (" + obj.Name + ") ended in " + (Ok ? "Success" : "Failure"));
                resolve(Ok);
            });
        });

        // Handle errors
        req.on('error', (error) => {
            console.error('Error during test execution:', error.message);
            Ok = false;
            console.log("Test (" + obj.Name + ") ended in " + (Ok ? "Success" : "Failure"));
            resolve(Ok);
        });

        // If the method is POST, write the content to the request body
        if (obj.Method === "POST") {
            req.write(postData);
        }

        // End the request
        req.end();
    });
}


async function Main() {
    let Ok = true; // Initialize Ok variable to true
    console.log("////////////////");

    try {
        const data = await fs.readFile(filePath, 'utf8');
        const jsonData = JSON.parse(data);

        for (const obj of jsonData) {
            if (! await RunTest(obj)) {
                Ok = false;
            }
        }
    } catch (error) {
        console.error('Error:', error);
        Ok = false;
    }

    console.log("////////////////");
    Ok ? console.log("Testing ended in Success") : console.log("Testing ended in Failure");
}

Main();