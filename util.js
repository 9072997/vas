function serializeArrayBuffer(buf) {
	let binary = '';
	let bytes = new Uint8Array(buf);
	let len = bytes.byteLength;
	for (var i = 0; i < len; i++) {
		binary += String.fromCharCode(bytes[i]);
	}
	return window.btoa(binary);
};

function deserializeArrayBuffer(str) {
	let binary_string = window.atob(str);
	let len = binary_string.length;
	let bytes = new Uint8Array(len);
	for (var i = 0; i < len; i++) {
		bytes[i] = binary_string.charCodeAt(i);
	}
	return bytes.buffer;
};
