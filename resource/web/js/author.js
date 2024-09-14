function openLoginModal() {
    $('#login').modal('show');
    return false;
}

function submitLoginModal() {
    document.getElementById('loginForm').submit();
    $('#login').modal('hide');
}


function openRegisterModal() {
    $('#register').modal('show');
    return false;
}

function submitRegisterModal() {
    document.getElementById('registerForm').submit();
    $('#register').modal('hide');
}

// retain login or not login status
function retainStatus() {
    const uid = getCookie('uid');
    const uname = getCookie('uname');
    if (uname.length > 0) {
        $('#status').html(`Hello, [${uid}] ${uname}!`);
        console.log(uid, uname);
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    else return '';
}

document.addEventListener('DOMContentLoaded', () => {
    retainStatus();
})
