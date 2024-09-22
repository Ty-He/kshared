async function fetchNotify() {
    const notify_ul = document.getElementById('notify-items');
    notify_ul.innerHTML = ''; // clear old data
    const res = await fetch('/get_unread_notify');
    if (!res.ok) {
        alert('network error');
        return;
    }

    const notifies = await res.json();

    notifies.forEach(notify => {
        const notify_li = document.createElement('li');
        notify_li.oncontextmenu = readThisNotify;
        notify_li.className = 'media';
        notify_li.setAttribute('notify_id', notify.id);


        const left = document.createElement('div');
        left.className = 'media-left';
        left.innerHTML = `<div class="media-object"></div>`;

        const body = document.createElement('div');
        body.className = 'media-body';
        
        const heading = document.createElement('div');
        heading.className = 'midia-heading';
        heading.innerHTML = `<span>&gt; receive: ${notify.sender} ${notify.release}</span>`;

        const content = document.createElement('span');
        // TODO text-trunc
        content.className = 'text-trunc';
        content.textContent = notify.content;

        const target = document.createElement('a');
        target.href = `/article?id=${notify.article_id}`;
        target.textContent = 'goto';

        body.appendChild(heading);
        body.appendChild(content);
        body.appendChild(target);

        notify_li.appendChild(left);
        notify_li.appendChild(body);

        notify_ul.appendChild(notify_li);
    });
}


async function readThisNotify(event) {
    event.preventDefault();
    const notify_id = event.currentTarget.getAttribute('notify_id');

    const res = await fetch(`/marked_notify_read?notify_id=${notify_id}`, {
        method: 'POST',
    });

    if (!res.ok) {
        alert('network error');
    } else {
        alert('You have marked it as read.');
    }
}
