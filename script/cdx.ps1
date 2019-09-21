
function cdx(
) {
    begin{
        [System.Collections.Generic.List[string]]$paths=@{}
    }
    process{
        $paths.Add($_)
    }
    end{
        $command = "$($paths | go-cdx $args)";
        Invoke-Expression "$command"
    }
}
