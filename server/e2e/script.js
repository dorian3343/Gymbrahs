const fs = require('fs').promises;
const https = require('https');
const http = require('http');

const filePath = 'tests.json';
async function RunTest(obj) {

    let errMessage = "Malformed Test: ";

    if (obj.Url === undefined) {
        console.log("Test (" + obj.Name + ") failed -> Malformed Test -> Missing URL");
        return false;
    }

    if (obj.Method === undefined) {
        console.log("Test (" + obj.Name + ") failed -> Malformed Test -> Missing HTTP method");
        return false;
    }

    try {
        let response;
        if (obj.Content !== undefined) {
            response = await fetch(obj.Url, {
                method: obj.Method,
                body: JSON.stringify(obj.Content),
                headers: { 'Content-Type': 'application/json' }
            });
        } else {
            response = await fetch(obj.Url, { method: obj.Method });
        }

        const json = await response.json();
        if (JSON.stringify(json) !== JSON.stringify(obj.ExpectedContent) && obj.ExpectedContent != undefined){
            console.log("Test (" + obj.Name + ") failed -> Got unexpected response -> " + json);
            return false
        }

        if (obj.Expected != response.status){
            console.log("Test (" + obj.Name + ") failed -> Got unexpected status -> " + response.status);
            return false
        }
    } catch (error) {
        console.log("Test (" + obj.Name + ") failed -> " + error);
        return false;
    }
    console.log("Test (" + obj.Name + ") Succeeded");
    return true;
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