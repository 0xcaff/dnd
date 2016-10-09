// The real meaning of DnD is drag and drop.
var dropzone = document.getElementById("dropzone");
var form = document.getElementById("form");
var uploads = document.getElementById("uploads");

function addEventsListener(element, events, func, prop) {
	events.forEach((e) => element.addEventListener(e, func, prop));
}

addEventsListener(dropzone, ['dragover' ,'dragenter'], (e) => {
	e.preventDefault();

	dropzone.classList.add('dragover');
});

addEventsListener(dropzone, ['dragleave' ,'dragend', 'drop'], (e) => {
	e.preventDefault();

	dropzone.classList.remove('dragover');
});

addEventsListener(dropzone, ['drop'], (e) => {
	handleFiles(e.dataTransfer.files);
});

addEventsListener(form, ['change'], (e) => {
	handleFiles(e.target.files);
});

function handleFiles(files) {
	dropzone.classList.add('uploading');

	var i = 0;
	for (let file of files) {
		let form = new FormData();
		form.set("file", file, file.name);

		let status = document.createElement("div");
		status.id = i;
		status.textContent = `${file.name} - 0%`;
		uploads.appendChild(status);

		let xhr = new XMLHttpRequest();
		xhr.open('POST', '/send');
		xhr.onload = () => console.log("Status: ", xhr.status, xhr.responseText);
		xhr.onerror = () => console.error("XHR Error");
		xhr.upload.onprogress = (event) => {
			if (event.lengthComputable) {
				var percent = Math.round((event.loaded / event.total) * 100);
				let old = status.textContent.split(" - ");
				old[old.length - 1] = `${percent}%`
				status.textContent = old.join(" - ");
			}
		}
		xhr.send(form);
	}
}

