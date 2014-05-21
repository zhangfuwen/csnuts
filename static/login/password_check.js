
function if_email_valid() {
	var email=document.getElementById("email").value;
	if(email.search(/^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/) == -1 ) {
		document.getElementById("email_valid").src="";
		document.getElementById("email_valid").alt="Email Invalid.";
		return false;
	}
	document.getElementById("email_valid").src="static/login/accept.png";
	return true;
}
function if_password_valid() {
	var passWord=document.getElementById("password").value;
	if( passWord.length<8 || passWord.length>16 ) {
//		document.getElementById("password_valid").src="static/login/reject.png";
		document.getElementById("password_valid").src="";
		document.getElementById("password_valid").alt="Password must be between 8~16 characters.";
		return false;
	}
	document.getElementById("password_valid").src="static/login/accept.png";
	if( document.getElementById("password1").value!="") {
		if_password_consistent();
	}
//	document.getElementById("password_valid").alt="OK";
	return true;
}
function if_password_consistent() {
	var passWord=document.getElementById("password").value;
	var passWord1=document.getElementById("password1").value;
	if( passWord != passWord1) {
//		document.getElementById("accept").src="static/login/reject.png";
		document.getElementById("accept").src="";
		document.getElementById("accept").alt="The passwords you entered is different.";
		return false;
	}
	document.getElementById("accept").src="static/login/accept.png";
//	document.getElementById("accept").alt="OK";
	return true;
}
