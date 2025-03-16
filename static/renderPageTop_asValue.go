package static

var RenderPageTop = `
 <div class="container sticky-top bg-dark-subtle mt-0 w-100 rounded-5 rounded-top-0">
 <div class="row">
 <div class="col-md-auto align-self-center text-light">
 <div class="bg-success rounded px-2">
 		<p class="text-center fw-bold"> App name <p>
		<p class="text-center"> {{.AppName}} <p>
		</div>
	</div>

	<div class="col-md-auto">
	<div class="bg-info rounded px-2 text-light">
		<p class="text-center"> Total memory </p>
		<p class="dropdown-divider"></p>
		<p class="text-center"> {{.Total}}  ->  {{.TotalKB}}KB</p>
		</div>
	</div>
 </div>
 </div>

 <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>

`
