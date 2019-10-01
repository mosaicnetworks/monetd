.. _monetd_systemd_rst:

Monetd Systemd Service
----------------------

Here is an example service file defining a ``systemd`` service for monetd:

``cat /etc/systemd/system/monet.service``:

.. code::

    [Unit]
    Description=monet node
    Requires=network-online.target
    
    [Service]
    User=admin
    ProtectSystem=strict
    NoNewPrivileges=yes
    PrivateTmp=yes
    PrivateDevices=yes
    DevicePolicy=closed
    ProtectHome=yes
    ProtectControlGroups=yes
    ProtectKernelModules=yes
    ProtectKernelTunables=yes
    RestrictAddressFamilies=AF_INET AF_INET6
    RestrictRealtime=yes
    RestrictNamespaces=yes
    MemoryDenyWriteExecute=yes
    Restart=on-failure
    RestartSec=3
    LimitNOFILE=32768
    ReadWritePaths=/opt/monet/data
    ExecStart=/opt/monet/bin/monetd run -d /opt/monet/data
    
    [Install]
    WantedBy=multi-user.target

It is fairly locked down and prevents from writing outside of 
``/opt/monet/data``.

Note that this requires ``monetd`` to be installed in ``/opt/monet/bin`` and for
the configuration to have been initialised in ``/opt/monet/data``. Here, we run 
the service as the ``admin`` user, which should have enough permissions in those
directories.

You can then use ``systemctl`` and ``journalctl`` to start, stop, and monitor
the monetd daemon:

.. code:: 

    systemctl start monet # start monetd
    journalctl --unit=monet # logs
    sytstectl stop monet # stop monetd
  
