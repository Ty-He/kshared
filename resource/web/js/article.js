
// when article loading, fecth md file from server.
document.addEventListener('DOMContentLoaded', async () => {
    let query = $('#content').text();
    document.getElementById('content').innerHTML = '';

    document.getElementById('content').classList.add('boxLoading');
    document.getElementById('content').classList.remove('invisible');

    let res = await fetch('/article_content?' + query);
    if (!res.ok) {
        alert(`HTTP error! Status: ${res.status}`);
        return;
    }

    let content = await res.text();
    document.getElementById('content').classList.remove('boxLoading');
    document.getElementById('content').innerHTML = marked.parse(content); 
    generate_toc(document.getElementById('content'), document.getElementById('toc'));
});

// search all <h> in content, 
// then generate toc and set to dest
function generate_toc(content, dest) {
    const headings = content.querySelectorAll('h1, h2, h3');
    const toc = document.createElement('ul');
    toc.className = 'toc';
    let cur_level = 1;
    // a stack for ul 
    let stack = [toc];

    headings.forEach(heading => {
        const level = parseInt(heading.tagName.substring(1));
        const toc_item = document.createElement('li');
        switch (level) {
            case 1:
                toc_item.className = "level-1";
                break;
            case 2:
                toc_item.className = "level-2";
                break;
            case 3:
                toc_item.className = "level-3";
                break;
        }
        const a = document.createElement('a');
        if (!heading.id) {
            heading.id = heading.textContent.toLowerCase().replace(/\s+/g, '-');
        }
        a.href = `#${heading.id}`;
        a.textContent = heading.textContent;
        a.addEventListener('mouseover', function() {
            this.style.color = 'black';
            this.style.textDecoration = 'none';
        });
        a.addEventListener('mouseout', function() {
            this.style.color = ''; 
        });


        toc_item.appendChild(a);

        if (level > cur_level) {
            // h1 h2 [h3 h4]
            const new_toc = document.createElement('ul');
            new_toc.className = 'toc';
            stack[stack.length - 1].appendChild(new_toc);
            stack.push(new_toc);
        } else if (level < cur_level) {
            while (cur_level > level) {
                stack.pop();
                cur_level--;
            }
        }

        stack[stack.length - 1].appendChild(toc_item);
        cur_level = level;

    });

    dest.innerHTML = '';
    dest.appendChild(toc);
}
