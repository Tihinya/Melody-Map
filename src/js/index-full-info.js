import "./googleApi.js"

let map
let data

export async function getParams(id) {
  const data = await getData(id)
  data.zoom = 1

  for (let v of data) {
    new google.maps.Marker({
      position: v,
      map: map,
    })
  }

  map.setCenter(data[0])
  map.setZoom(10)
}

async function getData(id) {
  const response = await fetch("/dateslocations/" + id)
  const json = await response.json()
  return json
}

export function centerLocation(id) {
  if (id >= data.lenght) {
    return
  }
  map.setCenter(data[0])
}

window.initMap = () => {
  map = new google.maps.Map(document.getElementById("map"))
}
