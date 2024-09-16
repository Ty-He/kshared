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
        $('#status').html(`Hello, [${uid}]${uname}!`);
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

// upload md flie

function openUploadModal() {
    $('#upload').modal('show');
    return false;
}

function submitUploadModal() {
    document.getElementById('uploadForm').submit();
    $('#upload').modal('hide');
}

// const uploadForm = document.getElementById('uploadForm');

// uploadForm.addEventListener('submit', (event) => {
//     event.preventDefault();
//     // author is assign 
//     const formData = new FormData(uploadForm);
//     const author = formData.get('author');
//     // curUser is form cookie
//     const uid = getCookie('uid');
//     if (uid.length == 0) {
//         alert('Please login with your id and pwd.');
//         uploadForm.reset();
//         return;
//     }
//     const uname = getCookie('uname');
//     const curUser = `${uid}&${uname}`;
//     if (author == null) {
//         // use defalut value
//         document.querySelector('input[name="author"]').value = curUser;
//     } else if (author != curUser) {
//         alert('If you want to write, please ensure correct.');
//         uploadForm.reset();
//         return;
//     }

//     uploadForm.submit();
//     uploadForm.reset(); // clear
// });

// update file

function openUpdateModal() {
    $('#update').modal('show');
    return false;
}

// request: update a article
async function updateArticle(formData) {
    // query = '' if lack query str
    let query = window.location.search; // include '?'
    let res = await fetch(`/update${query}`, {
        method: 'POST',
        body: formData,
    });
    if (!res.ok) {
        alert('Bad Requst!');
    } else {
        document.getElementById('home').click();
        alert('Update article ok!');
    }
}

function submitUpdateModal() {
    // document.getElementById('updateForm').submit();
    // override form submit:
    updateArticle(new FormData(document.getElementById('updateForm')));
    document.getElementById('updateForm').reset();  
    $('#update').modal('hide');
}

// delete file
async function deleteCurrentArticle() {
    let query = window.location.search; // include '?'
    let res = await fetch(`/delete${query}`, {
        method: 'DELETE',
    });
    if (!res.ok) {
        alert('Bad Requst!');
    } else {
        alert('Delete article ok!');
        // return home
        document.getElementById('home').click();
    }
}
