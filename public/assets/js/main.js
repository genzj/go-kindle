(function () {
    Dropzone.autoDiscover = false;

    var myDropzone = new Dropzone(
        '#my-dropzone',
        {
            autoProcessQueue: false,
            addRemoveLinks: true,
        }
    );

    document.getElementById('upload').addEventListener(
        'click',
        () => myDropzone.processQueue()
    );

    function getCookie(name) {
        var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
        if (match) return match[2];
    }

    myDropzone.on(
        'sending',
        function (file, xhr, formData) {
            formData.append('_csrf', getCookie('_csrf'));
        }
    )
    myDropzone.on(
        'success',
        (file) => {
            myDropzone.processQueue();
            setTimeout(
                () => myDropzone.removeFile(file),
                700
            );
        }
    )
})()
