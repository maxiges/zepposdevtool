package static

var RenderMainPage = `
<html><body>

<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
<div class="container text-center my-4">
    <h1 class="display-5 fw-bold text-primary mb-2">
        Zepp Metric Measurement App
    </h1>
    
    <p class="lead text-secondary mb-4">
        Welcome! Here is the list of available applications and metrics:
    </p>
</div>
`

var RenderAppButtonPage = `
<div class="container text-center">
    <div class="row align-items-center justify-content-center bg-secondary-subtle rounded p-2">
        <div class="col-3 align-items-center ">
            <span>{{.AppName}}</span>
        </div>    
        <div class="col-3 align-items-center">    
            <a href="{{.AppName}}" class="btn btn-primary btn-lg" role="button" >Show</a>
        </div>
    </div>
</div>
`

var RenderNoAppFound = `
<div class="container my-5">
    <div class="row justify-content-center">
        <div class="col-lg-8 col-md-10">
            <div class="alert alert-warning text-center p-4 border-2 border-warning" role="alert">
                <svg xmlns="http://www.w3.org/2000/svg" width="60" height="60" fill="currentColor" class="bi bi-exclamation-triangle-fill mb-3" viewBox="0 0 16 16">
                    <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.91 1.024L8.432 10h-.863L7.09 6.024c-.044-.562.375-1.024.91-1.024zm.002 6a1 1 0 1 0 0 2 1 1 0 0 0 0-2z"/>
                </svg>
                <h2 class="alert-heading text-danger fw-bold">
                    No application found
                </h2>
                <p class="lead mb-4">
                   Please make sure you sent data or reload the page
                </p>
            </div>
        </div>
    </div>
</div>
`
