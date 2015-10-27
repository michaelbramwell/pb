function initialize() {
  var latlong = new google.maps.LatLng(-31.94603, 115.83546);
  var mapOptions = {
    zoom: 16,
    center: latlong,
    mapTypeId: google.maps.MapTypeId.ROADMAP
  }
  var map = new google.maps.Map(document.getElementById("map_canvas"), mapOptions);

  var marker = new google.maps.Marker({
	    position: latlong,
	    map: map,
	    title: ''
	});
}

function loadScript() {
  var script = document.createElement("script");
  script.type = "text/javascript";
  script.src = "http://maps.googleapis.com/maps/api/js?key=AIzaSyDfBKqKqgooh60KhddEQKZjWzth_U9dfFo&sensor=true&callback=initialize";
  document.body.appendChild(script);
}

window.onload = loadScript;

