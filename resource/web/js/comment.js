/*
* comment(a db record)
  * id(uint32)
  * pid(uint32, if null, is root)
  * article_id(uint32)
  * sender_id(uint32)
  * content(string)
  * release(time.Time)
*/

var current_sending_btn;
function openSendingModal(event) 
{
    if (isLogin()) {
        $('#sending_modal').modal('show');
        current_sending_btn = event.target;
    } else {
        alert('当前未登录！');
    }
    event.preventDefault(); 
    // event.stopPropagation()
}

// pid, article_id, contenter
// id->increasing key; release->time.Now(); sender_id->cookie
async function sendingComment(event) 
{
    $('#sending_modal').modal('hide');
    const content = document.getElementById('reply_content').value;
    const pid = CommentElement(current_sending_btn).getAttribute('comment-id');
    const article_id = window.location.search.substring(4); // ?id=

    const data = {
        pid: pid,
        article_id, article_id,
        content, content,
    };
    console.log(data);

    try {
        const res = await fetch('/sending_comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data),
        });
        if (!res.ok) {
            throw new Error('Bad HTTP Reuqest');
        }
        
        // sending ok
        alert('Sending comment is ok!');
    } catch (error) {
        // sending bad
        alert(error);
    }
    // 
    document.getElementById('reply-content').value = '';
    event.preventDefault(); 
}



// send get request for comment, server should response json.
/*
comment object
{ 
    id:         key
    sender:     the name of sender
    release:       the tmie of release
    content:    ...
    target:     recver name
}
*/

async function fetchComment(event) {
    const article_id = window.location.search.substring(4); 
    // current selected comment
    const id = CommentElement(event.target).getAttribute('comment-id');
    const query = `?id=${id}&article_id=${article_id}`;
    
    const res = await fetch('/fetch_comment' + query);
    if (!res.ok) {
        alert('Bad Request!');
        return;
    } 

    // let data = await res.json();
    // let comments = JSON.parse(data);
    const comments = await res.json();
    if (comments.length == 0) {
        return;
    }
    // create a <ul> after span
    const comments_ul = document.createElement('ul');
    comments_ul.className = 'media-list';
    comments.forEach(comment => {
        const comment_li = document.createElement('li');
        comment_li.className = 'media';
        comment_li.setAttribute('comment-id', comment.id);
        
        // left
        const left = document.createElement('div');
        left.className = 'media-left';
        left.innerHTML = `<div class="media-object">${comment.sender}</div>`;

        // body
        const body = document.createElement('div');
        body.className = 'media-body';
        
        const heading = document.createElement('div');
        heading.className = 'media-heading';
        heading.innerHTML = `
            <span>&gt; reply ${comment.target} ${comment.release}</span>
            <div class="media-buttons">
                <a href="#" onclick="openSendingModal(event)"><span class="glyphicon glyphicon-comment"></span></a>
                &emsp13;
                <a herf="#" onclick="fetchComment(event)"><span class="glyphicon glyphicon-triangle-bottom"></span></a>
            </div>`;
        const content = document.createElement('span');
        content.innerHTML = marked.parse(comment.content);

        body.appendChild(heading);
        body.appendChild(content);

        comment_li.appendChild(left);
        comment_li.appendChild(body);

        comments_ul.appendChild(comment_li);
    });
    // <media-body>
    const nextCommentParent = event.target.parentElement.parentElement.parentElement.parentElement;

    const old_ul = nextCommentParent.querySelector('ul.media-list');
    if (old_ul) nextCommentParent.removeChild(old_ul);

    nextCommentParent.appendChild(comments_ul);

    // call once => bad: may be call more than one for the new comment
    // event.target.parentElement.removeAttribute('onclick');
    
    event.preventDefault(); 
}

function CommentElement(cur_elem) {
    return cur_elem.parentElement.parentElement.parentElement.parentElement.parentElement;
}

function isLogin() {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; uid=`);
    return parts.length == 2;
}
