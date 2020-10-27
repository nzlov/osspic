package main

const (
	_lt = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="zh-CN" lang="zh-CN">
<body>
    <h2>
        <ul>
            {{range .}} 
                <li><a href="{{.Path}}">{{.Name}}</a></li>
            {{end}}  
        </ul>  
    </h2>
</body>
</html>`
	_it = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="zh-CN" lang="zh-CN">
<meta http-equiv="Pragma" content="no-cache">
<meta http-equiv="Cache-Control" content="no-cache">
<meta http-equiv="Expires" content="0">
<body>
    <h2>
    <div>
        <a href="/list/{{.Bucket}}/{{.Name}}">跳过</a>
    </div>
    <div>
        <a href="/list/{{.Bucket}}/{{.Rame}}/270">左转</a>
        <a href="/list/{{.Bucket}}/{{.Rame}}/90">右转</a>
        <a href="/list/{{.Bucket}}/{{.Rame}}/180">180度</a>
    </div>
    <div>
        <a href="/list/{{.Bucket}}/{{.Rame}}/l270">本地左转</a>
        <a href="/list/{{.Bucket}}/{{.Rame}}/l90">本地右转</a>
        <a href="/list/{{.Bucket}}/{{.Rame}}/l180">本地180度</a>
    </div>
    <div>
        <img style="width: 640px;height:480px" src="{{.Path}}"></img>
        <p>{{.Path}}</p>
    </div>
    </h2>
</body>
</html>`
)
