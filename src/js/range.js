function showValue1(newValue) {
  document.getElementById("output1").innerHTML = newValue;
}

function showValue2(newValue) {
  document.getElementById("output2").innerHTML = newValue;
}

function showValue3(newValue) {
  document.getElementById("output3").innerHTML = newValue;
}

function showValue4(newValue) {
  document.getElementById("output4").innerHTML = newValue;
}





document.getElementById("range-from-cd").oninput = (e) => {
  showValue1(e.target.value)
}


document.getElementById("range-to-cd").oninput = (e) => {
  showValue2(e.target.value)
}


document.getElementById("range-from-fa").oninput = (e) => {
  showValue3(e.target.value)
}


document.getElementById("range-to-fa").oninput = (e) => {
  showValue4(e.target.value)
}