// The real meaning of DnD is drag and drop.

var dropzone = document.body;

var form = document.getElementById("form");
var uploads = document.getElementById("uploads");
var infoTemplate = document.getElementById("info");
var errors = document.getElementById("error");

// Handle UI Updates
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
		// Create Form
		let form = new FormData();
		form.set("file", file, file.name);

		// Display Status Info
		let localInfoTemplate = document.importNode(infoTemplate.content, true);
		localInfoTemplate.id = i;

		let progress = localInfoTemplate.getElementById("progress");
		let status = localInfoTemplate.getElementById("details");
		status.textContent = `${file.name} - 0%`;

		uploads.appendChild(localInfoTemplate);

		// Prepare Request
		let xhr = new XMLHttpRequest();
		xhr.addEventListener('load', (event) => {
			console.log("Status: ", xhr.status, xhr.responseText);
			progress.classList.remove("active");
		});

		xhr.addEventListener('error', (event) => {
			progress.classList.add("error");
			status.textContent += " - Something went wrong.";
		});

		xhr.upload.addEventListener('progress', (event) => {
			// TODO: Download Speed

			let old = status.textContent.split(" - ");
			if (event.lengthComputable) {
				// Update Progress
				var progressString = `${Math.round((event.loaded / event.total) * 100)}%`;
				old[1] = progressString;
				progress.style.setProperty("width", progressString);
			} else {
				// Remove Progress
				old.splice(1, 1);
			}
			status.textContent = old.join(" - ");
		});

		// Send Request
		xhr.open('POST', '/send');
		xhr.send(form);
	}
}

