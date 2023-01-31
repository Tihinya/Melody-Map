// Initialize and add the map
async function initMap() {
  const data = await getData();
  const map = new google.maps.Map(document.getElementById("map"), {
    zoom: 4,
    center: data[0],
  });

  for (let i = 0; i < data.length; i++) {
    new google.maps.Marker({
      position: data[i],
      map: map,
    });
  }
}

async function getData() {
  const response = await fetch("/dateslocations");
  const json = await response.json();
  return json;
}

window.initMap = initMap;
