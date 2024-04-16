async function Refresh() {
    const req = await fetch('/auth/refresh', {
        method: 'POST',
    });

    if (req.status !== 200) {
        //200番以外
        return false;
    }

    //更新確定
    const sreq = await fetch('/auth/refreshs', {
        method: 'POST',
    })

    if (sreq.status !== 200) {
        //200番以外
        return false;
    }

    return true;
}

Refresh();