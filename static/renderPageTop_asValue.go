package static

var RenderMainPage = `
<html><body>

<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
<div mb-2>
<p class="text-center">Welcome to the zepp metric measurement app</p>
<p class="text-center">Here is the list of available applications:</p>
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
