/*const acl = new Accelerometer({ frequency: 60 });
acl.addEventListener("reading", () => {
    document.querySelector('.acs_info_1').textContent = acl.x;
    document.querySelector('.acs_info_2').textContent = acl.y;
    document.querySelector('.acs_info_3').textContent = acl.z;
});

acl.start();

let gyroscope = new Gyroscope({frequency: 60});

gyroscope.addEventListener('reading', e => {
    document.querySelector('.gyr_info_1').textContent = gyroscope.x;
    document.querySelector('.gyr_info_2').textContent = gyroscope.y;
    document.querySelector('.gyr_info_3').textContent = gyroscope.z;
});
gyroscope.start();

document.querySelector('.gyr_info_3').textContent = 'huiii';
const handleOrientation = (event) => {
    const absolute = event.absolute;
    document.querySelector('.gyr_info_3').textContent = event.alpha;
    document.querySelector('.gyr_info_1').textContent = event.beta;
    document.querySelector('.gyr_info_2').textContent = event.gamma;
    //...
};

window.addEventListener("ondeviceorientation", handleOrientation, true);*


function orientation(event){
console.log("Magnetometer: "
    + event.alpha + ", "
    + event.beta + ", "
    + event.gamma
);
}
if(window.DeviceOrientationEvent){
window.addEventListener("deviceorientation", orientation, false);
}else{
console.log("DeviceOrientationEvent is not supported");
}*/




var absolute = false;
var rot_x = 0;
var rot_y = 0;
var rot_z = 0;
var acc_x = 0;
var acc_y = 0;
var acc_z = 0;
var ball_x = 150.5;
var ball_y = 150.5;
var move_x = 0;
var move_y = 0;
var looping = true;
var gravity = false;
var ball = null;

$(document).ready(function() {
	ball = document.getElementById('ball');

	window.addEventListener('deviceorientation', function( event ) {
		absolute = event.absolute;
        // X -180 - 180
        rot_x = event.beta.toFixed(2);
        // Y -90 - 90, loops twice over
        rot_y = event.gamma.toFixed(2);
        // Z 0 - 360
        rot_z = event.alpha.toFixed(2);
		
	}, false);

	function handleMotionEvent(event) {
		// with gravity 9,807 m/s²
        acc_x = event.accelerationIncludingGravity.x.toFixed(2);
        acc_y = event.accelerationIncludingGravity.y.toFixed(2);
        acc_z = event.accelerationIncludingGravity.z.toFixed(2);
    
	}
	window.addEventListener("devicemotion", handleMotionEvent, true);

	var currentScreenOrientation = window.orientation || 0;
	window.addEventListener('orientationchange', function() {
		currentScreenOrientation = window.orientation;
	}, false);

	window.setInterval(function(){
		if (looping)
		{
			$("#data").text(``);
			$("#data").append(`
			Absolute (z 0 is north): ${absolute}
			<br>Rotation X ${rot_x}
			<br>Rotation Y ${rot_y}
			<br>Rotation Z ${rot_z}
			<br>Acceleration X ${acc_x}
			<br>Acceleration Y ${acc_y}
			<br>Acceleration Z ${acc_z}
			<br>screen rotated ${currentScreenOrientation} degrees`);
			
		} else {
			$("#data").text(`Phone and browser must support motion sensors and they must be allowed for this site.
			 Hold the phone screen up and press Start. Rotate your phone to check sensor data. 
			 Visual representation uses x and y axis, rotation is clamped within 30 degrees.`);
		}
	}, 10);
});

function clamp(num, min, max) {
	return num <= min ? min : num >= max ? max : num;
}