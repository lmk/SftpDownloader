/*
별도의 html 파일 없이 실행 파일에 포함하기 위한 코드
*/
package main

import "fmt"

const HTML_ROOT = `
<!DOCTYPE html>
<html lang="kr">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.12.9/dist/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>  
	<title>Sftp Downloader</title>
</head>
<body class="bg-light">
	<div class="container">
		<div class="py-5 text-center">
			<h2>Sftp Downloader</h2>
			<p class="lead">경로를 유지하고 파일 전체를 다운로드 받습니다</p>
		</div>
		<div class="col-md-12">
			<form action="http://localhost:%s/download" method="post">
				<h4 class="mb-3">SFTP Server</h4>
				<div class="row mb-1">
					<div class="input-group col-md-6 mb-3">
						<label class="input-group-text" for="sftp-addr">IP Address</label>
						<input type="text" class="form-control" name="sftp-addr" placeholder="" value="%s" required>
					</div>
				</div>
				<div class="row mb-3">
					<div class="input-group col-md-3 mb-3">
						<label class="input-group-text" for="sftp-id">ID</label>
						<input type="text" class="form-control" name="sftp-id" placeholder="" value="%s" required>
					</div>
					<div class="input-group col-md-3 mb-3">
						<label class="input-group-text" for="sftp-pwd">Password</label>
						<input type="text" class="form-control" name="sftp-pwd" placeholder="" value="%s" required>
					</div>
				</div>
				<h4 class="mb-3">Local</h4>
				<div class="row mb-3">
					<div class="input-group col-md-12 mb-3">
						<label class="input-group-text" for="local-dir">Local Directory</label>
						<input type="text" class="form-control" name="local-dir" placeholder="" value="%s" required>
					</div>
				</div>
				<h4 class="mb-3">File List</h4>
				<div class="row mb-3">
					<div class="col-md-12">
						<textarea class="col-md-12" rows="10" name="file-list" placeholder="File List...">%s</textarea>
					</div>
				</div>
				<hr class="mb-3">
				<button class="btn btn-primary btn-lg btn-block" type="submit">Download</button>
			</form>
		</div>
	</div>
</body>
</html>
`

const HTML_DOWNLOAD = `
<!DOCTYPE html>
<html lang="kr">
<head>
    <meta charset="utf-8">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.2.1.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/popper.js@1.12.9/dist/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>  
    <title>Sftp Downloader</title>
</head>
<body class="bg-light">
    <div class="container">
        <div class="py-5 text-center">
            <h2>Sftp Downloader</h2>
            <p class="lead">경로를 유지하고 파일 전체를 다운로드 받습니다</p>
        </div>
        <div class="col-md-12">
            <table class="table">
                <thead>
                    <tr>
                    <th scope="col">#</th>
                    <th scope="col">Stat</th>
                    <th scope="col">Path</th>
                    <th scope="col">Date</th>
                    <th scope="col">Size</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
                </table>
        </div>
    </div>
	<script>
		var x = setInterval(function() {
			$.ajax({
				type: 'post',
				url: 'http://localhost:%s/downloading',
				success : function(result) {

					if (result.stat != "DOWNLOADING") {
						console.log(result)
						clearInterval(x)
						return
					} 

					$.Each(result.files, function(index, el){
						ch = ""
						per = el.localSize / el.remoteSize * 100
						if (per <= 0) {
							ch = '&#128347;'
						} else if (per <= 25) {
							ch = '&#128338;'
						} else if (per <= 50) {
							ch = '&#128341;'
						} else if (per <= 75) {
							ch = '&#128344;'
						} else {
							ch = '✔'
						}
						
						$("#" + el.path).text = ch
					})
				},
				error : function(xhr, status, message) {
					clearInterval(x)
					console.log("error : "+message);
					window.open('','_self').close(); 
				}
			})
			if ( true ) {
				clearInterval(x)
			}
		}, 1000);
	</script>
</body>
</html>
`

const HTML_DOWNLOAD_ROW = `
<tr>
	<th scope="row">%d</th>
	<td><div id="%s">%s</div></td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
</tr>
`

func HtmlRoot() string {
	return fmt.Sprintf(HTML_ROOT,
		AppPort,
		cfg.Ip,
		cfg.Id,
		cfg.Password,
		cfg.LocalDir,
		FilesPathToString(cfg.RemoteFiles))
}

func HtmlDownload() string {

	html := ""
	for i, file := range cfg.RemoteFiles {
		if file.isExists {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, i+1, file.path, "✔", file.path, file.date, HumanSize(float64(file.size)))
		} else {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, i+1, file.path, "❌", file.path, file.date, HumanSize(float64(file.size)))
		}

	}

	return fmt.Sprintf(HTML_DOWNLOAD, html, AppPort)
}

func HtmlDownloading() string {

	html := ""
	for i, file := range cfg.LocalFiles {

		if file.isExists {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, i+1, file.path, "✔", file.path, file.date, HumanSize(float64(file.size)))
		} else {
			html += fmt.Sprintf(HTML_DOWNLOAD_ROW, i+1, file.path, "❌", file.path, file.date, HumanSize(float64(file.size)))
		}

	}

	return fmt.Sprintf(HTML_DOWNLOAD, html, AppPort)
}
