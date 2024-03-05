import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import {ConvertFiles} from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';

var time = 0;
var percent = 0
EventsOn('app:took', (data) => {

    const [uuid, prgressPercent,filename] = data.split(',');

    CreteNewProgressBar(uuid,prgressPercent,filename);

    const value = parseInt(prgressPercent);
    var eleme = document.getElementById(uuid);

    eleme.value = value;
});

function CreteNewProgressBar(uuid , startValue,fileName) {

    if (!document.getElementById(uuid) || document.getElementById(uuid).id !== uuid) {
        const divElement = document.createElement('div');
        divElement.className = "container";

        const file = document.createElement('p');
        file.innerText = fileName;

        const UUID = document.createElement('p');
        UUID.innerText = uuid;

        const progressBar = document.createElement("progress");
        progressBar.id = uuid;
        progressBar.max = 100;
        progressBar.value = startValue;
        divElement.appendChild(file);
        divElement.appendChild(UUID);
        divElement.appendChild(progressBar);

        document.body.appendChild(divElement);

    }
};

EventsOn('app:tick', (data) => {
    console.log(data)
    console.log("here");
    const time = data; // Assuming 'data' is defined elsewhere

    console.log(time);
    let timeElement = document.getElementById("1");

    if (!timeElement) {
        timeElement = document.createElement("text");
        timeElement.id = "1";
        timeElement.innerText = time;
        document.body.appendChild(timeElement);
    } else {
        timeElement.innerText = time;
    }
});


document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
      <div class="result" id="result">👇enter Input and outPath below 👇</div>
      <div class="input-box" id="input">
        <input class="input" id="inputPathfield" type="text" placeholder="Input Path Here" autocomplete="off" />
         <input class="input" id="outputpathfield" type="text" placeholder="Output Path Here" autocomplete="off" />
         <button class="btn" onclick="convertFiles()">Convert</button>
         <div class = pro id="pro">
      </div>
    </div>
`;
document.getElementById('logo').src = logo;


let inp = document.getElementById("inputPathfield");
let outp = document.getElementById("outputpathfield");

window.convertFiles = function () {


    let inputPath =  inp.value;
    let outputpath = outp.value;
    console.log()
    // Check if the input is empty
   // if (inputPath === "") return;
  //  if (outputpath === "") return;
    // Call App.Greet(name)
    try {
        ConvertFiles(inputPath,outputpath)
            .then((result) => {
                // Update result with data back from App.Greet()
                console.log("in the convert try stament . then ")
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};
